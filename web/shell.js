
const ws = new WebSocket('ws://localhost:1111/ws');
ws.onopen = function(){
    console.log('Connection established');
};

ws.onmessage = function(event) {
    const Data = prettyJSON(event.data)
    $('#data').html(`<pre><span>${Data}</span></pre>`);
    $('#data').fadeIn(500);

    //$('body').animate({scrollTop:0}, 4000);

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
        let res = JSON.stringify(jsonObject, null, 4);
        return  res
    } catch (error) {
        console.log("invalid json")
        return jsonString;
  }
}

