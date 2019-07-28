package com.cp.callgov.api

import android.widget.Toast
import com.cp.callgov.model.Document
import com.google.gson.Gson
import io.reactivex.Single
import okhttp3.Call
import okhttp3.MediaType
import okhttp3.MediaType.Companion.toMediaTypeOrNull
import okhttp3.OkHttpClient
import okhttp3.Request
import retrofit2.Response
import java.io.IOException
import java.util.concurrent.TimeUnit

object ApiService {
    private val api = Api.create()
    var token = ""
    val client = OkHttpClient.Builder()
        .readTimeout(3, TimeUnit.SECONDS)
        .build()

    fun login(login: String, password: String): Single<Response<LoginResponse>> {
        val request = LoginRequest(login, password)
        return api.login(request)
    }

    fun createDocument(d: Document) = Single.create<String> { emitter ->
        val str = Gson().toJson(d)
        val body = okhttp3.RequestBody.create("application/json; charset=utf-8".toMediaTypeOrNull(), str)
        val request = Request.Builder()
            .url("http://45.128.204.157/json/obr?token=$token")
            .post(body)
            .build()
        val call = client.newCall(request)
        call.enqueue(object : okhttp3.Callback {
            override fun onFailure(call: Call, e: IOException) {
                emitter.onError(e)
            }

            override fun onResponse(call: Call, response: okhttp3.Response) {

                emitter.onSuccess(response.body!!.string())
            }


        })

    }

    fun fetchDocuments() = Single.create<String> { emitter ->
        val request = Request.Builder()
            .url("http://45.128.204.157/json/obrlist?token=$token")
            .get()
            .build()
        val call = client.newCall(request)
        call.enqueue(object : okhttp3.Callback {
            override fun onFailure(call: Call, e: IOException) {
                emitter.onError(e)
            }

            override fun onResponse(call: Call, response: okhttp3.Response) {

                emitter.onSuccess(response.body!!.string())
            }


        })

    }

}