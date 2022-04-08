package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strconv"

	"example.com/m/product-api/data"
)

type Products struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProduct(rw, r)
		return
	}
	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}
	if r.Method == http.MethodPut {
		p.l.Println("PUT", r.URL.Path)
		//expect the id in URL
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)
		if len(g) != 1 {
			p.l.Println("Invalid URI more than one id")
			http.Error(rw, "Invalid URL - 1", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 2 {
			p.l.Println("Invalid URI more than one capture group")
			http.Error(rw, "Invalid URL - 2", http.StatusBadRequest)
		}
		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			p.l.Println("Invalid URI unable to convert number", idString)
			http.Error(rw, "Invalid URL-3", http.StatusBadRequest)
			return
		}
		p.l.Println("got id", id)

		// various ways to find typeOf of the variable
		/*p.l.Printf("got id %T", id)
		p.l.Println("got id", reflect.TypeOf(id))
		p.l.Println("got id", reflect.ValueOf(id).Kind())*/

		p.updateProduct(id, rw, r)
		return

	}
	// catch all
	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) getProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET products")
	lp := data.GetProduct()
	/*d, err := json.Marshal(lp)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}*/
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to encode JSON", http.StatusInternalServerError)
	}
	//rw.Write(d)
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Post Product")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}

	p.l.Printf("Prod: %#v", prod)
	data.AddProduct(prod)
}

func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle PUT method: Inside here")
	prod := &data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
	}
	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product Not Found - 1", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "Product Not Found - 2", http.StatusInternalServerError)
		return
	}
}
