package raft

import (
	"encoding/binary"
	"fmt"
	"github.com/AllenShaw19/raft/log"
	"google.golang.org/protobuf/proto"
	"os"
	"syscall"
)

type ProtoBufFile struct {
	path string
}

func (f *ProtoBufFile) Save(message proto.Message, sync bool) error {
	tmpPath := f.path + ".tmp"
	file, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|syscall.FD_CLOEXEC, 0666)
	if err != nil {
		log.Error("open file failed, path: %s, err %v", f.path, err)
		return err
	}
	defer file.Close()

	headerBuf := make([]byte, 0, 8)
	msgBuf, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	binary.LittleEndian.PutUint32(headerBuf, uint32(len(msgBuf)))

	if n, err := file.Write(headerBuf); err != nil || n != 8 {
		log.Error("write len failed, path: %s, written len %d, err %v", tmpPath, n, err)
		return fmt.Errorf("write header fail %v", err)
	}
	if n, err := file.Write(msgBuf); err != nil || n != len(msgBuf) {
		log.Error("write failed, path: %s, written len %d, err %v", tmpPath, n, err)
		return fmt.Errorf("write msg fail %v", err)
	}

	if sync {
		err := file.Sync()
		if err != nil {
			log.Error("sync failed, path: %s", tmpPath)
			return err
		}
	}

	err = os.Rename(tmpPath, f.path)
	if err != nil {
		log.Error("rename failed, old: %s, new: %s", tmpPath, f.path)
		return err
	}

	return nil
}

func (f *ProtoBufFile) Load(message proto.Message) error {
	file, err := os.OpenFile(f.path, os.O_RDONLY, 0666)
	if err != nil {
		log.Error("open file failed, path: %s, err %v", f.path, err)
		return err
	}

	defer file.Close()

	headerBuf := make([]byte, 8)
	if n, err := file.Read(headerBuf); err != nil || n != 8 {
		log.Error("read len failed, path: %s, read len: %d, err: %v", f.path, n, err)
		return fmt.Errorf("%v", err)
	}
	len := binary.LittleEndian.Uint32(headerBuf)

	msgBuf := make([]byte, len)
	if n, err := file.ReadAt(msgBuf, 8); err != nil || uint32(n) != len {
		log.Error("read body failed, path: %s, read len: %d, err: %v", f.path, n, err)
		return fmt.Errorf("%v", err)
	}

	err = proto.Unmarshal(msgBuf, message)
	if err != nil {
		return err
	}

	return nil
}
