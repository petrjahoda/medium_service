const testButton = document.getElementById("test-button")
const buttonTime = document.getElementById("button-time")

testButton.addEventListener("click", function () {
    buttonTime.textContent = "Button clicked at: " + new Date().toLocaleString("en-IE");
})


const time = new EventSource('/time');
console.log(time)
time.addEventListener('time', (e) => {
    document.getElementById("actual-time").innerHTML = "Actual time using SSE: " + e.data;

}, false);