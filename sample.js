var a = [1,2,3,4,5]

var b = [5,6,7]

function c(one,two){
	
	return one+two;
}


console("script started")

setTimeout(
	function(configuredTimeout) {
		console("timeout called with configured delay = "+configuredTimeout);
	}, 3 * 1000
)

setInterval(
	function(counter) {
		console("counter = "+counter);
	}, 1000
)

run(
	function(answer){
		console("answer: "+ answer);
	}, "123"
)

console("script finished")

JSON.stringify( {"a":a, "b":b, "c":c(5,6)})
