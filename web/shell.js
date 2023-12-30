// hestory queriers
// queryID is a identity of entir saved queries
var queryID = 0
// WebSocket
function Connection() {
var ws = new WebSocket('ws://localhost:1111/ws');
ws.onopen = function(){
    console.log('Connection established');
}


ws.onerror = function(){
    console.log('Connection error');
    setTimeout(function() {
        Connection();
    }, 3000) // 3 second
}


ws.onmessage = function(event) {
    const Data = prettyJSON(event.data)
    $('#examples').hide();
    $('#data').html(`<pre><span>${Data}</span></pre>`);
    $('#data').fadeIn(500);
};

//  when ws closed reconnect after 2 second
ws.onclose = function() {
    console.log('WebSocket connection closed');
    $('#reconnecte').show();
    console.log("reconnet after 3 second")
    setTimeout(function() {
        Connection()
    }, 3000) // 3 second
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

Connection()

function saveQuery() {
            // save query-val in localstor for history
            console.log($('textarea').val())
            localStorage.setItem("{queryID}", textareaValue);
            queryID++
}

