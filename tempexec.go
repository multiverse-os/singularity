package singularity

//import (
//	"archive/zip"
//	"io"
//	"io/ioutil"
//	"log"
//	"os"
//	"os/exec"
//	"path/filepath"
//)

//func execEmbeddedBinary(tmpDir string, f *zip.File) ([]byte, error) {
//	tmpName := filepath.Join(tmpDir, f.Name)
//	rc, err := f.Open()
//	if err != nil {
//		return nil, err
//	}
//	defer rc.Close()
//	tmpFile, err := os.OpenFile(tmpName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
//	if err != nil {
//		return nil, err
//	}
//	_, err = io.Copy(tmpFile, rc)
//	if err != nil {
//		return nil, err
//	}
//	err = tmpFile.Close()
//	if err != nil {
//		return nil, err
//	}
//	return exec.Command(tmpName).CombinedOutput()
//}
