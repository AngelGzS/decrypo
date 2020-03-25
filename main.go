package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ajdnik/decrypo/build"
	"github.com/ajdnik/decrypo/decryptor"
	"github.com/ajdnik/decrypo/file"
	"github.com/ajdnik/decrypo/pluralsight"
	"github.com/cheggaaa/pb/v3"
)

func main() {
	defClip, err := pluralsight.GetClipPath()
	if err != nil {
		panic(err)
	}
	defDb, err := pluralsight.GetDbPath()
	if err != nil {
		panic(err)
	}
	clips := flag.String("clips", defClip, "location of clip .psv files")
	db := flag.String("db", defDb, "location of sqlite file")
	output := flag.String("output", "./Pluralsight Courses/", "location of decrypted courses")
	version := flag.Bool("v", false, "print tool version")
	flag.Parse()

	if *version {
		fmt.Println(build.Version())
		os.Exit(0)
	}

	courses := pluralsight.CourseRepository{
		Path: *db,
	}
	clipCount, err := courses.ClipCount()
	if err != nil {
		panic(err)
	}
	svc := decryptor.Service{
		Decoder: &pluralsight.Decoder{},
		Storage: &file.Storage{
			Path: *output,
		},
		CaptionEncoder: &file.SrtEncoder{},
		Clips: &pluralsight.ClipRepository{
			Path: *clips,
		},
		Courses: &courses,
	}
	bar := pb.StartNew(clipCount)
	err = svc.DecryptAll(func(_ decryptor.Clip, _ *string) {
		bar.Increment()
	})
	bar.Finish()
	if err != nil {
		panic(err)
	}
}
