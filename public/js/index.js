$(document).ready(function () {
    var options = {
        dataType: 'json',
        beforeSubmit: showSpinner,
        success: processJSON,
        error: showError
    };
    $('#subscribeUser').ajaxForm();
    $('#subscribeUser').submit(function () {
        $(this).ajaxSubmit(options); 
        return false;
    });
});

function showSpinner() {
    $(".spinner").show();
}

function showError() {
    var message = "There seems to be an error with the server. Please come back later."
    $(".spinner").hide();
    $(".message").text(message);
    $(".message").show();
}

function processJSON(data) {
    if (data.success === true) {
        $(".message").hide();
        $(".form").hide();
    }
    setTimeout(function () {
        $(".spinner").hide();
        $(".message").text(data.message);
        $(".message").show();
    }, 1000);
}