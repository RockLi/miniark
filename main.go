package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"

	"archive/tar"
	"compress/gzip"

	minio "github.com/minio/minio-go"
)

const (
	miniArkVersion  = "v1.0"
	minikubeVersion = "v0.32.0"
	k8sVersion      = "v1.12.4"
)

const (
	miniArkHome      = ".miniark"
	miniArkCacheHome = ".miniark/cache"
)

const (
	minioEndPoint        = "minio.longguikeji.com"
	minioAccessKey       = "public-ak"
	minioSecretAccessKey = "public-sk"
)

const (
	minikubeCmd = "minikube"
)

func Untar(dst string, r io.Reader) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
}

func getMiniArkHome() string {
	user, _ := user.Current()
	return fmt.Sprintf("%s/%s", user.HomeDir, miniArkHome)
}

func getMiniArkCacheHome() string {
	user, _ := user.Current()
	return fmt.Sprintf("%s/%s/%s", user.HomeDir, miniArkCacheHome, miniArkVersion)
}

func getMinioHome() string {
	return fmt.Sprintf("%s", miniArkVersion)
}

func getMinikubeHome() string {
	user, _ := user.Current()
	return fmt.Sprintf("%s/.minikube", user.HomeDir)
}

func downloadFromMinio(key, dest string) {
	// @Todo: add md5 checking and show progress

	fmt.Printf("downloading %s and save to %s\n", key, dest)

	_, err := os.Stat(dest)
	if err == nil {
		fmt.Println("file exists, skip...")
		return
	}

	minioClient, err := minio.New(minioEndPoint, minioAccessKey, minioSecretAccessKey, true)
	if err != nil {
		panic(err)
	}

	reader, err := minioClient.GetObject("ark", key, minio.GetObjectOptions{})
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	localFile, err := os.Create(dest)
	if err != nil {
		panic(err)
	}
	defer localFile.Close()

	stat, err := reader.Stat()
	if err != nil {
		panic(err)
	}

	if _, err := io.CopyN(localFile, reader, stat.Size); err != nil {
		panic(err)
	}
}

func downloadFiles() {
	// Download k8s files
	key := fmt.Sprintf("%s/k8s/minikube-%s", getMinioHome(), minikubeVersion)
	localFile := fmt.Sprintf("%s/k8s/minikube-%s", getMiniArkCacheHome(), minikubeVersion)
	_, err := os.Stat(fmt.Sprintf("%s/k8s", getMiniArkCacheHome()))
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(fmt.Sprintf("%s/k8s", getMiniArkCacheHome()), os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	downloadFromMinio(key, localFile)

	key = fmt.Sprintf("%s/k8s/minikube-%s.iso", getMinioHome(), minikubeVersion)
	localFile = fmt.Sprintf("%s/k8s/minikube-%s.iso", getMiniArkCacheHome(), minikubeVersion)
	downloadFromMinio(key, localFile)

	key = fmt.Sprintf("%s/k8s/%s/kubeadm", getMinioHome(), k8sVersion)
	localFile = fmt.Sprintf("%s/k8s/%s/kubeadm", getMiniArkCacheHome(), k8sVersion)
	_, err = os.Stat(fmt.Sprintf("%s/k8s/%s", getMiniArkCacheHome(), k8sVersion))
	if err != nil && os.IsNotExist(err) {
		err := os.MkdirAll(fmt.Sprintf("%s/k8s/%s", getMiniArkCacheHome(), k8sVersion), os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	downloadFromMinio(key, localFile)

	key = fmt.Sprintf("%s/k8s/%s/kubelet", getMinioHome(), k8sVersion)
	localFile = fmt.Sprintf("%s/k8s/%s/kubelet", getMiniArkCacheHome(), k8sVersion)
	downloadFromMinio(key, localFile)

	key = fmt.Sprintf("%s/k8s/%s/kubectl", getMinioHome(), k8sVersion)
	localFile = fmt.Sprintf("%s/k8s/%s/kubectl", getMiniArkCacheHome(), k8sVersion)
	downloadFromMinio(key, localFile)

	// Download core images
	key = fmt.Sprintf("%s/k8s/%s/core-images.tar.gz", getMinioHome(), k8sVersion)
	localFile = fmt.Sprintf("%s/k8s/%s/core-images.tar.gz", getMiniArkCacheHome(), k8sVersion)
	downloadFromMinio(key, localFile)

}

func prepareWorkspace() {
	var (
		home      = getMiniArkHome()
		cacheHome = getMiniArkCacheHome()
	)

	_, err := os.Stat(home)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(home, os.ModePerm)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	_, err = os.Stat(cacheHome)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(cacheHome, os.ModePerm)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
}

func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = os.Stat(dst)
	if err == nil {
		fmt.Printf("file %s already exists, will override automatically ...\n", dst)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	if err != nil {
		panic(err)
	}

	buf := make([]byte, 1024)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}

func prepareMinikube() {
	log.Printf("Preparing minikube environment")

	var (
		home = getMinikubeHome()
		src  string
		dst  string
	)

	if _, err := os.Stat(getMinikubeHome()); err == nil {
		os.RemoveAll(home)
	}

	os.MkdirAll(fmt.Sprintf("%s/cache/iso", home), os.ModePerm)
	os.MkdirAll(fmt.Sprintf("%s/cache/%s", home, k8sVersion), os.ModePerm)

	src = fmt.Sprintf("%s/k8s/minikube-%s.iso", getMiniArkCacheHome(), minikubeVersion)
	dst = fmt.Sprintf("%s/cache/iso/minikube-%s.iso", home, minikubeVersion)
	copyFile(src, dst)

	src = fmt.Sprintf("%s/k8s/%s/kubeadm", getMiniArkCacheHome(), k8sVersion)
	dst = fmt.Sprintf("%s/cache/%s/kubeadm", home, k8sVersion)
	copyFile(src, dst)

	src = fmt.Sprintf("%s/k8s/%s/kubelet", getMiniArkCacheHome(), k8sVersion)
	dst = fmt.Sprintf("%s/cache/%s/kubelet", home, k8sVersion)
	copyFile(src, dst)

	// restore core images
	localFile := fmt.Sprintf("%s/k8s/%s/core-images.tar.gz", getMiniArkCacheHome(), k8sVersion)
	r, _ := os.Open(localFile)
	dst = fmt.Sprintf("%s/cache", home)
	Untar(dst, r)

	// prepare config/config.json
	os.MkdirAll(fmt.Sprintf("%s/config", home), os.ModePerm)
	configFile := fmt.Sprintf("%s/config/config.json", home)

	data := `
	{
		"WantReportError": true,
		"cache": {
		"gcr.io/k8s-minikube/storage-provisioner:v1.8.1": null,
			"k8s.gcr.io/coredns:1.2.2": null,
			"k8s.gcr.io/etcd:3.2.24": null,
			"k8s.gcr.io/kube-addon-manager:v8.6": null,
			"k8s.gcr.io/kube-apiserver:v1.12.4": null,
			"k8s.gcr.io/kube-controller-manager:v1.12.4": null,
			"k8s.gcr.io/kube-proxy:v1.12.4": null,
			"k8s.gcr.io/kube-scheduler:v1.12.4": null,
			"k8s.gcr.io/kubernetes-dashboard-amd64_v1.10.1": null,
			"k8s.gcr.io/pause:3.1": null
		}
	
	}
	`

	ioutil.WriteFile(configFile, []byte(data), os.ModePerm)
}

func startMinikube() {
	fmt.Println("Start minikube...")
	cmd := exec.Command(minikubeCmd, "start", "--kubernetes-version", k8sVersion)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	output, _ := cmd.Output()
	fmt.Println(string(output))

}

func installK8SCli() {
	file := fmt.Sprintf("%s/k8s/%s/kubectl", getMiniArkCacheHome(), k8sVersion)
	dstFile := "/usr/local/bin/kubectl"
	if err := copyFile(file, dstFile); err != nil {
		panic(err)
		panic("failed to copy kubectl")
	}

	os.Chmod(dstFile, 0754)
}

func installMinikube() {
	minikubeFile := fmt.Sprintf("%s/k8s/minikube-%s", getMiniArkCacheHome(), minikubeVersion)
	dstFile := "/usr/local/bin/minikube"
	copyFile(minikubeFile, dstFile)

	os.Chmod(dstFile, 0754)
}

func installHelmCli() {
}

func enableMinikubeAddon(addon string) error {
	cmd := exec.Command(minikubeCmd, "addons", "enable", addon)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("minikube addon %s enable failed", addon)
	}

	return nil
}

func enableMinikubeAddons() error {
	var addons = []string{
		"ingress",
		"dashboard",
		"metrics-server",
	}

	for _, addon := range addons {
		if err := enableMinikubeAddon(addon); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	log.Printf("miniark: %s, kubernetes: %s, minikube: %s\n", miniArkVersion, k8sVersion, minikubeVersion)

	// prepareWorkspace()

	// downloadFiles()

	// prepareMinikube()
	// installMinikube()
	// startMinikube()

	installK8SCli()

	// if err := enableMinikubeAddons(); err != nil {
	// 	log.Fatal(err)
	// }

	fmt.Println("Happy hacking! Any issues please report to rock<insfocus@gmail.com>")
}
