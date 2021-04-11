package easyquery

import (
	"bytes"
	"easyquery/tools/json_util"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin/binding"
	jsoniter "github.com/json-iterator/go"
)

var JsonContentType = []string{"application/json; charset=utf-8"}

var BaseJsoner = InitBaseJsoner()

func InitBaseJsoner() jsoniter.API {

	jsoner := jsoniter.ConfigDefault
	jsoner.RegisterExtension(&json_util.NamingStrategyExtension{
		DummyExtension: jsoniter.DummyExtension{},
		Translate:      json_util.InitialLower,
	})

	return jsoner
}

type BaseJsonRender struct {
	Data   interface{}
	jsoner jsoniter.API
}

func NewBaseJsonRender(data interface{}) BaseJsonRender {
	return BaseJsonRender{
		Data:   data,
		jsoner: BaseJsoner,
	}
}

func (r BaseJsonRender) Render(w http.ResponseWriter) (err error) {
	if err = r.WriteJSON(w, r.Data); err != nil {
		panic(err)
	}
	return
}

func (r BaseJsonRender) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, JsonContentType)
}

func (r BaseJsonRender) WriteJSON(w http.ResponseWriter, obj interface{}) error {
	writeContentType(w, JsonContentType)
	jsonBytes, err := r.jsoner.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

// EnableDecoderUseNumber is used to call the UseNumber method on the JSON
// Decoder instance. UseNumber causes the Decoder to unmarshal a number into an
// interface{} as a Number instead of as a float64.
var EnableDecoderUseNumber = false

// EnableDecoderDisallowUnknownFields is used to call the DisallowUnknownFields method
// on the JSON Decoder instance. DisallowUnknownFields causes the Decoder to
// return an error when the destination is a struct and the input contains object
// keys which do not match any non-ignored, exported fields in the destination.
var EnableDecoderDisallowUnknownFields = false

var BaseJsonBinding = jsonBinding{jsoner: BaseJsoner}

type jsonBinding struct {
	jsoner jsoniter.API
}

func (b jsonBinding) Name() string {
	return "json"
}

func (b jsonBinding) Bind(req *http.Request, obj interface{}) error {
	if req == nil || req.Body == nil {
		return fmt.Errorf("invalid request")
	}
	return b.decodeJSON(req.Body, obj)
}

func (b jsonBinding) BindBody(body []byte, obj interface{}) error {
	return b.decodeJSON(bytes.NewReader(body), obj)
}

func (b jsonBinding) decodeJSON(r io.Reader, obj interface{}) error {
	decoder := b.jsoner.NewDecoder(r)
	if EnableDecoderUseNumber {
		decoder.UseNumber()
	}
	if EnableDecoderDisallowUnknownFields {
		decoder.DisallowUnknownFields()
	}
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return b.validate(obj)
}

func (b jsonBinding) validate(obj interface{}) error {
	if binding.Validator == nil {
		return nil
	}
	return binding.Validator.ValidateStruct(obj)
}
