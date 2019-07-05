$(document).ready(function(){
    var user_name
    var final_connection

    $("#form_registro").submit(function(e){
        e.preventDefault();
        user_name = $("#user_name").val();
        $.ajax({
            type:"POST",
            url: "http://localhost:4545/validate",
            data:{
                "user_name": user_name
            },
            success: function(data){
                result(data)
            }
        })
    })

    function result(data) {
        
        obj = JSON.parse(data)
        if(obj.IsValid === true ){
            createConexion()
        }else{
            console.log("Intentalo de nuevo");
            
        }
    }   

    function createConexion(){
        $("#container").hide();
        $("#container_chat").show();

        var conexion = new WebSocket("ws://localhost:4545/chat/" + user_name);
        final_connection = conexion;
        conexion.onopen = function(ressponse){
            conexion.onmessage = function(ressponse){
                console.log("Nos envio algo")
                console.log(ressponse.data);
                val = $("#chat_area").val();
                $("#chat_area").val(val + "\n" + ressponse.data)

                
            }
        }
    }

    $("#form_message").on("submit", function(e){
        e.preventDefault();
        mensaje = $("#msg").val();
        final_connection.send(mensaje);
        $("#msg").val("")
    })
})