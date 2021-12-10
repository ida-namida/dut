package httpserver

import (
    "dut/service"
    "dut/view"
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/julienschmidt/httprouter"
)

type Server struct {
    httpPort        string
    decisionService service.DecisionService
}

func NewServer(httpPort string, decisionService service.DecisionService) Server {
    return Server{httpPort, decisionService}
}

func (s Server) Start() {
    router := httprouter.New()
    router.GET("/decisions/:name/form", s.getDecisionInputForm)
    router.POST("/decisions/:name/evaluate", s.evaluateDecision)

    http.ListenAndServe(":" + s.httpPort, router)
}

func (s Server) getDecisionInputForm(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
    request := view.DecisionFormRequest{
        Name: params.ByName("name"),
    }

    response, err := s.decisionService.GetDecisionInputForm(request)
    if err != nil {
        w.Header().Add("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(fmt.Sprintf(`{
            "error": "%s"
        }`, err.Error())))
        return
    }

    w.Header().Add("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(response)
}

func (s Server) evaluateDecision(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
    decoder := json.NewDecoder(r.Body)

    name := params.ByName("name")
    request := make(map[string]interface{})
    err := decoder.Decode(&request)
    if err != nil {
        w.Header().Add("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        w.Write([]byte(fmt.Sprintf(`{
            "error": "Cannot parse request into JSON - %s"
        }`, err.Error())))
        return
    }

    response, err := s.decisionService.EvaluateDecision(name, request)
    if err != nil {
        w.Header().Add("Content-Type", "application/json")
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(fmt.Sprintf(`{
            "error": "%s"
        }`, err.Error())))
        return
    }

    w.Header().Add("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(response)
}