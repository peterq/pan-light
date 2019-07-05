package cv

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"os"
	"path"
	"sync"
)

var xmlFile = "./data/haarcascades_cuda/haarcascade_frontalface_alt.xml"
var classifier = gocv.NewCascadeClassifier()

func init() {
	if !classifier.Load(xmlFile) {
		fmt.Printf("Error reading cascade file: %v\n", xmlFile)
		return
	}
}

func CV() {

	deviceID := 0

	// open webcam
	webcam, err := gocv.VideoCaptureDevice(int(deviceID))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	// open display window
	window := gocv.NewWindow("Face Detect")
	defer window.Close()

	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	// color for the rect when faces detected
	blue := color.RGBA{0, 0, 255, 0}

	fmt.Printf("start reading camera device: %v\n", deviceID)
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("cannot read device %d\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// detect faces
		rects := classifier.DetectMultiScale(img)
		fmt.Printf("found %d faces\n", len(rects))

		// draw a rectangle around each face on the original image,
		// along with text identifying as "Human"
		for _, r := range rects {
			gocv.Rectangle(&img, r, blue, 3)

			size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
			pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
			gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
		}

		// show the image in the window, and wait 1 millisecond
		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}

var l = new(sync.Mutex)

func CheckFace(imgPath string, savePath string) []image.Rectangle {
	l.Lock()
	defer l.Unlock()
	blue := color.RGBA{0, 0, 255, 0}

	// 打开图片
	img := gocv.IMRead(imgPath, -1)
	defer img.Close()
	// 进行识别
	rects := classifier.DetectMultiScale(img)

	if savePath != "" {
		// 在图片中标记
		for _, r := range rects {
			gocv.Rectangle(&img, r, blue, 3)

			size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
			pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
			gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
		}
		os.MkdirAll(path.Dir(savePath), os.ModePerm)
		gocv.IMWrite(savePath, img)
	}

	return rects
}
