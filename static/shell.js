
// WebSocket
var ws = new ReconnectingWebSocket('ws://localhost:1111/ws', null, {debug: true, reconnectInterval: 2000});

function Connection() {

ws.onopen = function(){
    console.log('Connection established');
    $('#reconnecte').fadeOut(500);
}

//  when ws closed reconnect after 2 second
ws.onclose = function(e) {
    console.log('WebSocket connection closed'), e.error;
    //ws.close()

    $('#reconnecte').show();
    console.log("reload page to reconnect")
  };

ws.onerror = function(e){
    console.log('Connection error', e.error);
//    ws.close()

    $('#reconnecte').show();
//    console.log("reload page to reconnect")

}


ws.onmessage = function(event) {
    const Data = prettyJSON(event.data)
    $('#examples').hide();
    $('#data').html(`<pre><span>${Data}</span></pre>`);
    $('#data').fadeIn(500);

};

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

            return;
        } 
    }
});

function prettyJSON(jsonString) {
    try {
        const jsonObject = JSON.parse(jsonString);
        let res = JSON.stringify(jsonObject, null, 3);
        return  res
    } catch (error) {
        console.log("invalid json")
        return jsonString;
  }
}
}

// Dealing with Textarea Height
function calcHeight(value) {
  let numberOfLineBreaks = (value.match(/\n/g) || []).length;
    if (numberOfLineBreaks < 3) {
        numberOfLineBreaks = 3
    }
    console.log("lines:",numberOfLineBreaks)
  // min-height + lines x line-height + padding + border
  let newHeight = 20 + numberOfLineBreaks * 20 + 12 + 2;
    
  return newHeight;
}

let textarea = document.querySelector("textarea");
textarea.addEventListener("keyup", () => {
  textarea.style.height = calcHeight(textarea.value) + "px";
    console.log("height",calcHeight(textarea.value))
});
	

Connection()

