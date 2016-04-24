package static
import (
	"html/template"
	"net/http"
	"github.com/julienschmidt/httprouter"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

var (
	templateIndex = template.Must(template.ParseFiles("html/index.html"))
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	c := appengine.NewContext(r)
	log.Infof(c, "processing static.Index")
	err := templateIndex.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
