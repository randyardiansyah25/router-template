package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"router-template/delivery"
	"strings"

	"github.com/gin-gonic/gin"
)

// ** Struct ini dibuat untuk wrapping gin.ResponseWriter, karena gin tidak secara otomatis menyimpan atau memberikan akses ke
// ** body response setelah ditulis ke client (berbeda dengan request yang memberikan akses di gin.context.Request.Body).
// ** Tambah property body *bytes.Buffer untuk clone dari response body.
// ** Karena responseBodyWriter implement dari interface gin.ResponseWriter (yang digunakan oleh gin context)
// ** maka gin akan memanggil Write(b []byte), yang mana implement disini menuliskan ke body *bytes.Buffer.
// ** Body ini yang digunakan untuk kebutuhan cetak response dari handler. Tanpa menyalin, body di gin context
// ** akan kosong, sehingga data(body) yang di kirim ke client saat response pun menjadi kosong.
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// ** Implement dari interface gin.ResponseWriter
func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// ** Middleware untuk handling log saat response
func ResponseLogger(ctx *gin.Context) {
	w := &responseBodyWriter{
		body:           bytes.NewBufferString(""), //* inisiasi dulu, karena body ini pointer defaultnya nil
		ResponseWriter: ctx.Writer,                //* inisiasi responseBodyWriter mengunakan gin.ResponseWriter dari gin context
	}

	ctx.Writer = w

	//** Go to next handler
	ctx.Next()

	//** Ini setelah handler sebelumnya sudah memberikan response
	var bufferInfo strings.Builder
	bufferInfo.WriteString(fmt.Sprintf("[RESPONSE] %s, Path: %s, Response Info: \n", ctx.Request.Method, ctx.Request.URL.String()))
	bufferInfo.WriteString(fmt.Sprintf("Status Code: %d\n", w.Status()))

	//** Write Headers
	for name, values := range ctx.Writer.Header() {
		for _, value := range values {
			bufferInfo.WriteString(fmt.Sprintf("%s: %s\n", name, value))
		}
	}
	bufferInfo.WriteString("Body: ")
	bufferInfo.WriteString(w.body.String())

	respContentTypes := ctx.Writer.Header().Get("Content-Type")

	// Split content type dengan delimiter semicolon (;), karna umumnya header diisi beserta charset,
	// untuk pengecekan ini charset tidak dibutuhkan, maka buang saja charset nya. Contoh, Content-Type: application/json; charset=utf-8
	respContentType := strings.Split(respContentTypes, ";")[0]
	if respContentType == ContentTypeAppJson {
		var prettyJson bytes.Buffer
		if er := json.Indent(&prettyJson, w.body.Bytes(), "", "    "); er == nil {
			bufferInfo.WriteString(fmt.Sprintf("\n>> Beautify:\n%s\n", prettyJson.String()))
		}
	}

	delivery.PrintLog(bufferInfo.String())
}
