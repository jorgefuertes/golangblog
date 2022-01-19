package blog

import (
	"bufio"
	"golangblog/internal/log"
	"os"
	"path/filepath"
	"regexp"
)

func LoadAllPosts(root string) {
	path := root + "/post"
	log.Info("BLOG", "Loading all posts from: ", path)
	if err := iterateDir(path); err != nil {
		log.Warn("BLOG:WALK:ERROR", err)
	}
}

func iterateDir(path string) error {
	log.Debug("BLOG:WALK", "Walking into: ", path)
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		log.Debug("BLOG:LOOK", "Looking: ", info.Name())
		if info.Mode().IsRegular() && info.Name() == "post.md" {
			return loadPost(path)
		}
		return nil
	})
}

func loadPost(path string) error {
	log.Debug("BLOG:LOAD", "Loading post: ", path)
	p := new(Post)

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)
	mode := "none"
	for s.Scan() {
		if s.Text() == "[meta]" {
			mode = "meta"
			continue
		}

		if s.Text() == "[body]" {
			mode = "body"
			continue
		}

		if mode == "body" {
			p.Body += s.Text()
			continue
		}

		if mode == "meta" {
			r := regexp.MustCompile("^((?:[[:alnum:]]|[-_])+)[[:blank:]]*=[[:blank:]]*(.*)[[:blank:]]*$")
			matches := r.FindStringSubmatch(s.Text())
			if matches != nil {
				log.Debug("BLOG:META", matches[1], matches[2], " LEN ", len(matches))
			} else {
				log.Warn("BLOG:META", "Malformed? ", s.Text())
			}
		}
	}

	posts = append(posts, p)

	return nil
}
