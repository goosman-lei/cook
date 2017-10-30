package connector

import (
	"context"
	"errors"
	//"fmt"
	etcd "github.com/coreos/etcd/client"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrEtcdKeyNotFound = errors.New("key not found")
	ErrEtcdNotDir      = errors.New("not dir")
)

type EtcdConf struct {
	Addrs   []string
	Timeout time.Duration
}

type etcdWrapper struct {
	client etcd.KeysAPI
	Conf   EtcdConf
}

var (
	Etcd *etcdWrapper
)

func SetupEtcd(config EtcdConf) error {
	var err error
	Etcd, err = newEtcdWrapper(config)
	return err
}

func newEtcdWrapper(config EtcdConf) (*etcdWrapper, error) {
	e, err := etcd.New(etcd.Config{
		Endpoints: config.Addrs,
	})

	if err != nil {
		return nil, err
	} else {
		return &etcdWrapper{
			client: etcd.NewKeysAPI(e),
			Conf:   config,
		}, nil
	}
}

func (e *etcdWrapper) Mkdir(path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.Conf.Timeout)
	defer cancel()

	_, err := e.client.Set(
		ctx,
		path,
		"",
		&etcd.SetOptions{PrevExist: etcd.PrevIgnore, Dir: true})
	if err != nil && !strings.HasPrefix(err.Error(), "102: Not a file") {
		return err
	}

	return nil
}

func (e *etcdWrapper) Set(path, data string) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.Conf.Timeout)
	defer cancel()

	if _, err := e.client.Set(ctx, path, data, &etcd.SetOptions{PrevExist: etcd.PrevIgnore, Dir: false}); err != nil {
		if strings.HasPrefix(err.Error(), "100: Key not found") {
			dir := filepath.Dir(path)
			if err = e.Mkdir(dir); err == nil {
				_, err = e.client.Set(ctx, path, data, &etcd.SetOptions{PrevExist: etcd.PrevIgnore, Dir: false})
			}
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *etcdWrapper) Remove(path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), e.Conf.Timeout)
	defer cancel()

	_, err := e.client.Delete(
		ctx,
		path,
		&etcd.DeleteOptions{Recursive: true, PrevIndex: 0})
	if err != nil && strings.HasPrefix(err.Error(), "100: Key not found") {
		err = nil
	}

	return err
}

func (e *etcdWrapper) Get(path string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.Conf.Timeout)
	defer cancel()

	rsp, err := e.client.Get(
		ctx,
		path,
		&etcd.GetOptions{Recursive: false})
	if err != nil {
		if strings.HasPrefix(err.Error(), "100: Key not found") {
			return "", nil
		} else {
			return "", err
		}
	}

	return rsp.Node.Value, nil
}

func (e *etcdWrapper) Dir(path string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.Conf.Timeout)
	defer cancel()

	rsp, err := e.client.Get(
		ctx,
		path,
		&etcd.GetOptions{Recursive: false})
	if err != nil {
		if strings.HasPrefix(err.Error(), "100: Key not found") {
			return []string{}, nil
		} else {
			return nil, err
		}
	}

	if !rsp.Node.Dir {
		return nil, ErrEtcdNotDir
	}

	if len(rsp.Node.Nodes) < 1 {
		return []string{}, nil
	}

	items := []string{}
	for _, node := range rsp.Node.Nodes {
		items = append(items, node.Key)
	}

	return items, nil
}

func (e *etcdWrapper) IsDir(path string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), e.Conf.Timeout)
	defer cancel()

	rsp, err := e.client.Get(
		ctx,
		path,
		&etcd.GetOptions{Recursive: false})
	if err != nil {
		return false, err
	}

	return rsp.Node.Dir, nil
}
