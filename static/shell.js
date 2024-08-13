/*
// to fix sub field like "contact.phon",
for (let [key, value] of Object.entries(yourobject)) {
    console.log(key, value);
}
*/


// This configuration is suitable for development situation
//const configs = {debug: false, reconnectDecay:1 , reconnectInterval: 200, reconnectDecay:1, maxReconnectInterval:200}

// WebSocket
//var ws = new ReconnectingWebSocket('ws://localhost:1111/ws');

function connection() {
var ws = new WebSocket('ws://localhost:1111/ws');

ws.onopen = function(){
    console.log('Connection established');
    $('#reconnecte').fadeOut(500);
}

//  when ws closed reconnect after 2 second
ws.onclose = function() {
    ws.close()
    $('#reconnecte').show();
    setTimeout(connection(), 700);
  };

ws.onerror = function(){
    ws.close()
    //$('#reconnecte').show();
}

//
ws.onmessage = function(event) {
    const Data = prettyJSON(event.data)
    $('#examples').hide();
    $('#data').html(`<pre><span>${Data}</span></pre>`);
    $('#data').fadeIn(500);
};

//
const queryInput = document.getElementById('query-input');
queryInput.addEventListener('keydown', function(event) {
    if (event.altKey && event.key === 'Enter') {
        const cursorPosition = queryInput.selectionStart;

        // Insert a newline character at the cursor position
        const textBeforeCursor = queryInput.value.substring(0, cursorPosition);
        const textAfterCursor = queryInput.value.substring(cursorPosition);
        queryInput.value = textBeforeCursor + '\n' + textAfterCursor;

        // Move the cursor to the end of the newline
        queryInput.selectionStart = cursorPosition + 1;
        queryInput.selectionEnd = cursorPosition + 1;
        return 
    }
   
    if (event.key === 'Enter') {
        $("#data").css("display","none");
        event.preventDefault();
        if (queryInput.value) {
            eval("obj = "+ queryInput.value)
            let query = JSON.stringify(obj)
            ws.send(query);
            //console.log(query);
            return;
        } 
    }
});
} // end connection func

connection()


function prettyJSON(jsonString) {
     try {
        const jsonObject = JSON.parse(jsonString);
        let res = JSON.stringify(jsonObject, null, 3);
        return  res
     } catch (error) {
        console.error("invalid json")
        return jsonString;
    }
}


// Dealing with Textarea Height
function calcHeight(value) {

    let numberOfLineBreaks = (value.match(/\n/g) || []).length;
    if (numberOfLineBreaks < 3) {
        numberOfLineBreaks = 3
    }
    
    let text = textarea.value
    while (text.includes('\n\n\n'))  {
         text = text.replace("\n\n\n", "\n")
    }

    textarea.value = text 

    if (numberOfLineBreaks > 24) {
        numberOfLineBreaks = 24
    }


    //console.log("lines:",numberOfLineBreaks)
    // min-height + lines x line-height + padding + border
    let newHeight = 20 + numberOfLineBreaks * 20 + 12 + 2;
    
    return newHeight;
}

let textarea = document.querySelector("textarea");
textarea.addEventListener("keyup", () => {
  textarea.style.height = calcHeight(textarea.value) + "px";
    //console.log("height",calcHeight(textarea.value))
});
	


$(document).on("keypress", function (e) {
   // console.log("event : ", e) 
   // console.log("input : ", $("#query-input").val()) 
   // TODO some hilit for js object 
});


 
// copy paste example into textarea
$('pre').click(function () {
    $('textarea').val($(this).text())
    $('textarea').focus()
})

// pointer on pre examples
$('pre').css('cursor', 'pointer');

