var baseURL = "http://localhost:4000/";
var resource = "v1/summary";
var param = "?url="
var $ = function(id) { return document.getElementById(id) };

var form = $("search-form");
var input = $("search");
var results = $("results");
var searchIcon = $("search-icon");
form.addEventListener("submit", evt => {
    results.innerHTML = "";
    evt.preventDefault();   
    getMeta(evt);
});

searchIcon.onclick = function() {
    results.innerHTML = "";
    getMeta();
};

function getMeta(){        
    fetch(baseURL + resource + param + input.value)
        .then(function(response){
            if(response.ok){
                return response.json();                             
            }
        }).then(function(data){
            fillCard(data);
            console.log(data);
        }).catch(function(error) {
            console.log(error);
        });

}

function fillCard(data) {
    var card = document.createElement("div");
    card.classList.add("card");
    var content = document.createElement("div");
    content.classList.add("card-content");    
    var title = document.createElement("span");
    title.classList.add("card-title")
    title.innerText = data.title;
    var description = document.createElement("p");
    description.innerText = data.description;   
    if(data.videos != null) {
        //handle different video types
        var video = document.createElement("iframe")
        video.src = data.videos[0].url;
        video.classList.add("my-image");
        card.appendChild(video);
    } else if(data.images != null){
        var imageDiv = document.createElement("div");
        imageDiv.classList.add("card-image");
        var image = document.createElement("img");
        image.classList.add("my-image")
        image.src = data.images[0].url;
        imageDiv.appendChild(image);
        card.appendChild(image);
    }else if(data.icon != null){
        var image = document.createElement("img");
        image.classList.add("my-icon")
        image.src = data.icon.url;        
        card.appendChild(image);
    }
    content.appendChild(title);
    content.appendChild(description);    
    card.appendChild(content);
    results.appendChild(card);
}
