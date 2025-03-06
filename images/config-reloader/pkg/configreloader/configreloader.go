// Copyright © 2025 Kube logging authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configreloader

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/banzaicloud/config-reloader/pkg/metrics"
	fsnotify "github.com/fsnotify/fsnotify"
)

func New() (*ConfigReloader, error) {
	cfg := &ConfigReloader{}
	cfg.InitMode = flag.Bool("init-mode", false, "Init mode for unarchive files. Works only if volume-dir-archive exist. Default - false")
	cfg.DirForUnarchive = flag.String("dir-for-unarchive", "/tmp/unarchive", "Directory where the archives will be unpacked")
	cfg.Webhook.Method = flag.String("webhook-method", "POST", "the HTTP method url to use to send the webhook")
	cfg.Webhook.StatusCode = flag.Int("webhook-status-code", 200, "the HTTP status code indicating successful triggering of reload")
	cfg.Webhook.Retries = flag.Int("webhook-retries", 1, "the amount of times to retry the webhook reload request")

	flag.Var(&cfg.VolumeDirs, "volume-dir", "the config map volume directory to watch for updates; may be used multiple times")
	flag.Var(&cfg.VolumeDirsArchive, "volume-dir-archive", "the config map volume directory to watch for updates and unarchiving; may be used multiple times")
	flag.Var(&cfg.Webhook.Urls, "webhook-url", "the url to send a request to when the specified config map volume directory has been updated")
	flag.Parse()

	return cfg, nil
}

func (cfg *ConfigReloader) Run() error {

	err := cfg.checks()
	if err != nil {
		return err
	}

	if len(cfg.VolumeDirs) > 0 {
		err := cfg.volumeDirWatcher()
		if err != nil {
			return err
		}
	}
	if len(cfg.VolumeDirsArchive) > 0 {
		for _, vda := range cfg.VolumeDirsArchive {
			// fmt.Printf("VDA: %s\n", vda)
			cfg.unarchiveDir(vda)
		}
		if *cfg.InitMode {
			log.Println("Init mode completed")
			return nil
		}
		cfg.volumeDirArchiveWatcher()
	}

	return nil

}

func (cfg *ConfigReloader) volumeDirWatcher() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	// defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					continue
				}

				if !isValidEvent(event) {
					continue
				}

				log.Println("ConfigMap or Secret updated")
				cfg.sendWebHook()
			case err, ok := <-watcher.Errors:
				if !ok {
					continue
				}
				metrics.WatcherErrors.Inc()
				log.Println("Error:", err)
			}
		}
	}()

	for _, d := range cfg.VolumeDirs {
		log.Printf("Watching directory: %q", d)
		err = watcher.Add(d)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cfg *ConfigReloader) volumeDirArchiveWatcher() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	// defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					continue
				}

				if !isValidEvent(event) {
					continue
				}

				err := cfg.unarchiveDir(event.Name)
				if err != nil {
					log.Println("Error:", err)
				}

				log.Println("ConfigMap or Secret updated")
				cfg.sendWebHook()
			case err, ok := <-watcher.Errors:
				if !ok {
					continue
				}
				metrics.WatcherErrors.Inc()
				log.Println("Error:", err)
			}
			time.Sleep(time.Second * 10)
		}
	}()

	for _, d := range cfg.VolumeDirsArchive {
		log.Printf("Watching directory (with unarchive): %q", d)
		err = watcher.Add(d)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cfg *ConfigReloader) checks() error {
	if (len(cfg.VolumeDirs) < 1) && (len(cfg.VolumeDirsArchive) < 1) {
		return fmt.Errorf("%s", "Missing volume-dir or volume-dir-archive")
	}

	if len(cfg.Webhook.Urls) < 1 && !*cfg.InitMode {
		return fmt.Errorf("%s", "Missing webhook-url if initMode disable")
	}

	if *cfg.InitMode && (len(cfg.VolumeDirsArchive) < 1) {
		return fmt.Errorf("%s", "init-mode work only with volume-dir-archive")
	}

	if *cfg.InitMode && (len(cfg.VolumeDirs) > 0) {
		return fmt.Errorf("%s", "init-mode don't work with volume-dir")
	}

	return nil
}

func isValidEvent(event fsnotify.Event) bool {
	if event.Op&fsnotify.Create != fsnotify.Create {
		return false
	}

	if filepath.Base(event.Name) != "..data" {
		return false
	}
	return true
}

func (cfg *ConfigReloader) sendWebHook() {
	for _, h := range cfg.Webhook.Urls {
		begun := time.Now()
		req, err := http.NewRequest(*cfg.Webhook.Method, h.String(), nil)
		if err != nil {
			metrics.SetFailureMetrics(h.String(), "client_request_create")
			log.Println("Error:", err)
			continue
		}
		userInfo := h.User
		if userInfo != nil {
			if password, passwordSet := userInfo.Password(); passwordSet {
				req.SetBasicAuth(userInfo.Username(), password)
			}
		}

		successfulReloadWebhook := false

		for retries := *cfg.Webhook.Retries; retries != 0; retries-- {
			log.Printf("Performing webhook request (%d/%d)", retries, *cfg.Webhook.Retries)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				metrics.SetFailureMetrics(h.String(), "client_request_do")
				log.Println("Error:", err)
				time.Sleep(time.Second * 10)
				continue
			}
			resp.Body.Close()
			metrics.RequestsByStatusCode.WithLabelValues(h.String(), strconv.Itoa(resp.StatusCode)).Inc()
			if resp.StatusCode != *cfg.Webhook.StatusCode {
				metrics.SetFailureMetrics(h.String(), "client_response")
				log.Println("error:", "Received response code", resp.StatusCode, ", expected", cfg.Webhook.StatusCode)
				time.Sleep(time.Second * 10)
				continue
			}

			metrics.SetSuccessMetrics(h.String(), begun)
			log.Println("successfully triggered reload")
			successfulReloadWebhook = true
			break
		}

		if !successfulReloadWebhook {
			metrics.SetFailureMetrics(h.String(), "retries_exhausted")
			log.Println("error:", "Webhook reload retries exhausted")
		}
	}

}

func (cfg *ConfigReloader) unarchiveDir(path string) error {
	// fmt.Println(path)
	// kuberPath := path + "/..data"
	// fmt.Println(kuberPath)
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.Name()[len(file.Name())-3:] != ".gz" {
			continue
		}
		fullFilePath := path + "/" + file.Name()
		err := cfg.unarchiveFile(fullFilePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cfg *ConfigReloader) unarchiveFile(path string) error {
	outFileName := *cfg.DirForUnarchive + "/" + filepath.Base(path)[0:len(filepath.Base(path))-3]
	log.Printf("Unarhive file from %s to %s", path, outFileName)

	// if path[len(path)-3:] != ".gz" {
	// 	return fmt.Errorf("File %s is not a .gz archive. Do nothing", path)
	// }

	gzipFile, err := os.Open(path)
	if err != nil {
		return err
	}

	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	// Uncompress to a writer. We'll use a file writer
	outfileWriter, err := os.Create(outFileName)
	if err != nil {
		return err
	}
	defer outfileWriter.Close()

	// Copy contents of gzipped file to output file
	_, err = io.Copy(outfileWriter, gzipReader)
	if err != nil {
		return err
	}
	return nil
}
