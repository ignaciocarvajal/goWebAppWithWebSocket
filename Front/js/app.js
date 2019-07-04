$(document).ready(function(){

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
        
    }
})