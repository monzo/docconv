package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"

	"github.com/golang/protobuf/proto"

	"github.com/sajari/docconv/iWork"
	"github.com/sajari/docconv/snappy"
)

// Convert PAGES to text
func ConvertPages(r io.Reader) (string, map[string]string) {
	meta := make(map[string]string)
	var textBody string

	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println("ioutil.ReadAll:", err)
		return "", nil
	}

	zr, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		log.Println("zip.NewReader:", err)
		return "", nil
	}

	for _, f := range zr.File {
		if strings.HasSuffix(f.Name, "Preview.pdf") {
			// There is a preview PDF version we can use
			if rc, err := f.Open(); err == nil {
				return ConvertPDF(rc)
			}
		}
		if f.Name == "index.xml" {
			// There's an XML version we can use
			if rc, err := f.Open(); err == nil {
				return ConvertXML(rc)
			}
		}
		if f.Name == "Index/Document.iwa" {
			rc, _ := f.Open()
			defer rc.Close()
			bReader := bufio.NewReader(snappy.NewReader(io.MultiReader(strings.NewReader("\xff\x06\x00\x00sNaPpY"), rc)))
			archiveLength, err := binary.ReadVarint(bReader)
			archiveInfoData, err := ioutil.ReadAll(io.LimitReader(bReader, archiveLength))
			archiveInfo := &TSP.ArchiveInfo{}
			err = proto.Unmarshal(archiveInfoData, archiveInfo)
			fmt.Println("archiveInfo:", archiveInfo, err)
		}
	}

	return textBody, meta
}
