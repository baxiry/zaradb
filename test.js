(function rec(i){
    setTimeout(function(){

console.log(i)
        if(i <= 20) rec(i+1);
    }, 700);
})(0);

console.log("done")
