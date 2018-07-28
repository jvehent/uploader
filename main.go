package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	if os.Getenv("DESTDIR") == "" {
		log.Fatal("Set the DESTDIR env variable to the directory where uploaded files will be stored. eg. '/var/uploads/'")
	}
	if os.Getenv("BASEURL") == "" {
		log.Fatal("Set the BASEURL env variable to the base http url where uploaded files will be served from. eg 'https://example.net/files/'")
	}
	if os.Getenv("UPLOADURL") == "" {
		log.Fatal("Set the UPLOADURL env variable to the http url where upload requests are sent to. eg 'https://example.net/upload'")
	}
	http.HandleFunc("/", getIndex)
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(":5050", nil)
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s user-agent=%s", r.Method, r.Proto, r.URL.String(), r.UserAgent())
	w.Write(htmlHead)
	w.Write([]byte(fmt.Sprintf(`<body>
	<form id="uploadform" enctype="multipart/form-data" action=%q method="post">
	    <input type="file" name="uploadfile" />
	    <input type="submit" value="upload" />
	</form>
	<div class="progress">
	    <div class="bar"></div>
	    <div class="percent">0%%</div>
	</div>
</body>
</html>
`, os.Getenv("UPLOADURL"))))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s %s user-agent=%s", r.Method, r.Proto, r.URL.String(), r.UserAgent())
	if r.Method != http.MethodPost {
		http.Error(w, "only POST method is allowed", http.StatusBadRequest)
		return
	}
	log.Printf("parsing multipart form")
	// no more than 100MB of memory, the rest goes into /tmp
	r.ParseMultipartForm(100000000)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		http.Error(w, "failed to read uploadfile form field", http.StatusBadRequest)
		log.Printf("failed to read uploadfile form field: %v", err)
		return
	}
	defer file.Close()
	_, err = os.Stat(os.Getenv("DESTDIR") + handler.Filename)
	if err == nil {
		http.Error(w, fmt.Sprintf("a file named %q already exists at the destination", handler.Filename), http.StatusInternalServerError)
		log.Printf("a file named %q already exists at the destination", handler.Filename)
		return
	}
	go func() {
		log.Printf("writing file %v", handler.Header)
		f, err := os.OpenFile(os.Getenv("DESTDIR")+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			http.Error(w, "failed to write file to destination", http.StatusInternalServerError)
			log.Printf("failed to write file to destination: %v", err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}()
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("File uploaded successfully at " + os.Getenv("BASEURL") + handler.Filename))
}

var htmlHead = []byte(`<!DOCTYPE html>
<html>
<head>
       <title>Upload file</title>
	 <script src="https://code.jquery.com/jquery-2.2.4.min.js" integrity="sha256-BbhdlvQf/xTY9gja0Dq3HiwQF8LaCRTXxZKRutelT44=" crossorigin="anonymous"></script>
	 <script src="https://cdnjs.cloudflare.com/ajax/libs/jquery.form/4.2.2/jquery.form.min.js" integrity="sha384-FzT3vTVGXqf7wRfy8k4BiyzvbNfeYjK+frTVqZeNDFl8woCbF0CYG6g2fMEFFo/i" crossorigin="anonymous"></script>
	 <script>
	    $(function() {
		    var bar = $('.bar');
		    var percent = $('.percent');
		    $('form').ajaxForm({
			  beforeSend: function() {
			    var percentVal = '0%';
			    bar.width(percentVal);
			    percent.html(percentVal);
			  },
			  uploadProgress: function(event, position, total, percentComplete) {
			    var percentVal = percentComplete + '%';
			    bar.width(percentVal);
			    percent.html(percentVal);
			  },
			  complete: function(xhr) {
			    percent.html(xhr.responseText);
			    document.getElementById("uploadform").reset();
			    var uform = document.getElementById("uploadform");
			    uform.style.display = "none";
			  }
		    });
	    });
	</script>
</head>
`)
