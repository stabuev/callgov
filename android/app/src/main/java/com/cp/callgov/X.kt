package com.cp.callgov

import android.content.Context
import android.content.Intent
import com.cp.callgov.model.Document

fun Context.startDetail(d:Document){
    val i = Intent(this,DocumentCreationActivity::class.java)
    i.putExtra("obr",d)
    startActivity(i)
}