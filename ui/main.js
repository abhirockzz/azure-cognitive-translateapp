$(function() {
    //Translate text with flask route
    $("#translate").on("click", function(e) {
      e.preventDefault();
      var translateVal = document.getElementById("text-to-translate").value;
      var languageVal = document.getElementById("select-language").value;
      var translateRequest = { 'text': translateVal, 'to': languageVal }
  
      if (translateVal !== "") {
        $.ajax({
          url: '/translate-text',
          method: 'POST',
          headers: {
              'Content-Type':'application/json'
          },
          dataType: 'json',
          data: JSON.stringify(translateRequest),
          success: function(data) {
            for (var i = 0; i < data.length; i++) {
              document.getElementById("translation-result").textContent = data[i].translations[0].text;
              document.getElementById("detected-language-result").textContent = data[i].detectedLanguage.language;
              if (document.getElementById("detected-language-result").textContent !== ""){
                document.getElementById("detected-language").style.display = "block";
              }
              document.getElementById("confidence").textContent = data[i].detectedLanguage.score;
            }
          }
        });
      };
    });
  });