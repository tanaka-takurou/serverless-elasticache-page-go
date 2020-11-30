$(document).ready(function() {
  Get();
});

var Get = function() {
  const data = {action: "get"};
  request(data, (res)=>{
    $("#result").text(res.message);
  }, (e)=>{
    console.log(e.responseJSON.message);
    $("#warning").text(e.responseJSON.message).removeClass("hidden").addClass("visible");
  });
};

var Update = function() {
  const data = {action: "update"};
  request(data, (res)=>{
    $("#result").text(res.message);
  }, (e)=>{
    console.log(e.responseJSON.message);
    $("#warning").text(e.responseJSON.message).removeClass("hidden").addClass("visible");
  });
};

var request = function(data, callback, onerror) {
  $.ajax({
    type:          'POST',
    dataType:      'json',
    contentType:   'application/json',
    scriptCharset: 'utf-8',
    data:          JSON.stringify(data),
    url:           App.url
  })
  .done(function(res) {
    callback(res);
  })
  .fail(function(e) {
    onerror(e);
  });
};
var App = { url: location.origin + {{ .ApiPath }} };
