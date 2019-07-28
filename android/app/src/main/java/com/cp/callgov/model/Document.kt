package com.cp.callgov.model

import java.io.Serializable

data class Document(
    var id:Int = 0,
    var title:String = "",
    var content:String = "",
    var public:Int = 0,
    var state:String = "",
    var address:String = ""
):Serializable