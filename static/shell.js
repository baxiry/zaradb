


function connection() {
var ws = new WebSocket('ws://localhost:1111/ws');

ws.onopen = function(){
    console.log('Connection established');
    $('#reconnecte').fadeOut(400);
}

//  when ws closed reconnect after 700ms 
ws.onclose = function() {
    ws.close()
    $('#reconnecte').show();
    setTimeout(connection(), 700);
  };

ws.onerror = function(){
    ws.close()
}


// change json's output format 
var pretty = true 

$('#togglePretty').change(function() {
    pretty = !pretty;
    $('textarea').focus()
});

//
ws.onmessage = function(event) {
    if (pretty) {
        var Data = prettyJSON(event.data)
        $('#data').html(`<pre><span>${Data}</span></pre>`);
        $('#data').fadeIn(400);
        console.log("pretty json")
        return
    }


    $('#examples').hide();
    $('#data').html("<div><div>");

    Data = JSON.parse(event.data);;

    for (let i = 0;i< Data.length;i++) {
        let obj = JSON.stringify(Data[i]) 
        obj = obj.replace(/,"/g, ', "'); 

        $('#data').append(`<pre><span>${obj}</span></pre>`);
    }
    $('#data').fadeIn(400);
};

//
const queryInput = document.getElementById('query-input');
queryInput.addEventListener('keydown', function(event) {

    // 
    if ((event.metaKey || event.altKey ) && event.key === 'Enter' ) {
        $("#data").css("display","none");
        event.preventDefault();
        if (queryInput.value) {
            try {
                eval("obj = "+ queryInput.value)
                let query = JSON.stringify(obj)
                ws.send(query);
                return;

            } catch (error) {
                //console.error("Error evaluating code:", error);
                $('#data').html(`<pre><span>${error}</span></pre>`);
                $('#data').fadeIn(400);
                return; 
            }
        } 
    }

    if (event.key === 'Enter') {
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

})} // end connection func

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
    
    if (numberOfLineBreaks > 24) {
        numberOfLineBreaks = 24
    }

    return 20 + numberOfLineBreaks * 20 + 12;
}

let textarea = document.querySelector("textarea");

$(document).ready(function () {
   setHightTextArea(textarea)
})

textarea.addEventListener("keyup", (e) => {
   setHightTextArea(textarea)
});


function setHightTextArea(textarea) {
    let hi = calcHeight(textarea.value) 
    textarea.style.height = hi + "px";

    hi = calcHeight(textarea.value)  + 10 
    $("#output").css("height", "calc(100vh - "+ hi +"px)" )
    //css height: calc(100vh - 100px);
    //console.log("event : ", e) 
}

// copy paste example into textarea
$('pre').click(function () {
    $('textarea').val($(this).text())
    $('textarea').focus()
})

// pointer on pre examples
$('pre').css('cursor', 'pointer');

