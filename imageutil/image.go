package imageutil

import (
	"github.com/disintegration/imaging"
	"github.com/obnahsgnaw/goutils/errutil"
	"github.com/obnahsgnaw/goutils/pathutil"
	"image"
)

/*
Nearest最快但低质量Linear中等质量，速度较快Mitchell平滑边缘，适合缩小Lanczos最佳质量，但计算慢（推荐缩放）
*/

type Image struct {
	errutil.ErrBuilder
	filepath string
	src      image.Image
	dst      *image.NRGBA
}

func New(filepath string) (s *Image, err error) {
	s = &Image{
		filepath: filepath,
	}
	s.ErrPrefix = "image"

	s.filepath, err = pathutil.ValidFile(s.filepath)
	if err != nil {
		err = s.NewError(err, "image path invalid")
		s = nil
		return
	}

	s.src, err = imaging.Open(s.filepath)
	if err != nil {
		err = s.NewError(err, "open image[", s.filepath, "] failed")
		s = nil
		return
	}

	return
}

func (s *Image) Image() image.Image {
	return s.src
}

func (s *Image) Path() string {
	return s.filepath
}

func (s *Image) Thumbnail(w, h int) {
	s.dst = imaging.Thumbnail(s.img(), w, h, imaging.Lanczos)
}

func (s *Image) img() image.Image {
	if s.dst == nil {
		return s.src
	}
	return s.dst
}

func (s *Image) Save(dest string) error {
	if s.dst == nil {
		return imaging.Save(s.src, dest)
	}
	return imaging.Save(s.dst, dest)
}
