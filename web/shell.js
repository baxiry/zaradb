
const ws = new WebSocket('ws://localhost:1111/ws');
ws.onopen = function(){
    console.log('Connection established');
};


const dataOutput = document.getElementById('data');
const queryInput = document.getElementById('query-input');
console.log("queryInput is : ",queryInput)

ws.onmessage = function(event) {
    const Data = prettyJSON(event.data)
    //dataOutput.className = 'message';

    //$("#data").html(`<pre><span>${Data}</span></pre>`);
    dataOutput.innerHTML = `<pre><span>${Data}</span></pre>`;

    
};

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
        event.preventDefault();
        const message = queryInput.value;
        if (message) {
            ws.send(message);// message
            return;
        } 
    }
});

function prettyJSON(jsonString) {
    try {
        const jsonObject = JSON.parse(jsonString);
        let res = JSON.stringify(jsonObject, null, 4);
        //console.log(res)
        return  res
    } catch (error) {
        console.log("invalid json")
        return jsonString;
  }
}

