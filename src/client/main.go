package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"client/editor"

	"github.com/nsf/termbox-go"
)

var docId = flag.String("doc-id", "", "Document ID")

func main() {
	/* exporter, err := prometheus.New()
	if err != nil {
		fmt.Println("Error creating metrics exporter:", err)
	}
	meterProvider := otelmetric.NewMeterProvider(otelmetric.WithReader(exporter))
	traceExporter, err := otelstdouttrace.New(otelstdouttrace.WithPrettyPrint())
	if err != nil {
		fmt.Println("Error creating trace exporter:", err)
	}
	traceProvider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(traceExporter), sdktrace.WithResource(otelresource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName("grpc-client"))))
	textMapPropagator := otelpropagation.TraceContext{}
	dialOptions := opentelemetry.DialOption(opentelemetry.Options{
		MetricsOptions: opentelemetry.MetricsOptions{MeterProvider: meterProvider},
		TraceOptions:   oteltracing.TraceOptions{TracerProvider: traceProvider, TextMapPropagator: textMapPropagator},
	})

	go http.ListenAndServe("127.0.0.1:12345", promhttp.Handler()) */

	flag.Parse()

	err := termbox.Init()
	if err != nil {
		panic("Termbox init failed")
	}
	defer termbox.Close()

	// This input mode recognize escape characters
	termbox.SetInputMode(termbox.InputEsc)

	logFile, _ := os.OpenFile("editor" + string(RandomLetter()) + string(RandomLetter()) + string(RandomLetter()) + string(RandomLetter()) + string(RandomLetter()) + string(RandomLetter()) + string(RandomLetter()) + ".log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	log.SetOutput(logFile)
	defer logFile.Close()

	vEditor, recvUpdate := editor.NewEditor(*docId)

	shouldExit := false

	recvKeyDigit := make(chan *termbox.Event, 20)

	go func() {
		for !shouldExit {
			e := termbox.Event{
				Type: termbox.EventKey,
				Ch: RandomLetter(),
			}
			recvKeyDigit <- &e
			time.Sleep(time.Millisecond * time.Duration(RandomWaitTime()))
		}
	}()

	for !shouldExit {
		vEditor.Draw()
		select {
		case event := <-recvKeyDigit:
			{
				switch event.Type {
				case termbox.EventKey:
					shouldExit = vEditor.OnKeyEvent(*event)
				}
			}
		case <-recvUpdate:
			{
				fmt.Fprintf(editor.Writer, "Update from collaborators\n")
				editor.Writer.Flush()
			}
		}
	}
}


func RandomLetter() rune {
    if rand.Intn(2) == 0 {
        // Lowercase letter (a-z: ASCII 97-122)
        return rune(rand.Intn(26) + 97)
    } else {
        // Uppercase letter (A-Z: ASCII 65-90)
        return rune(rand.Intn(26) + 65)
    }
}

func RandomWaitTime() int {
    // Calculate range: 150-300 ms
    return rand.Intn(151) + rand.Intn(450)
}
