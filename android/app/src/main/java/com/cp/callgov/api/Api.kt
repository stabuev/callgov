package com.cp.callgov.api

import com.cp.callgov.model.Document
import io.reactivex.Observable
import io.reactivex.Single
import okhttp3.OkHttpClient
import okhttp3.logging.HttpLoggingInterceptor
import retrofit2.Response
import retrofit2.Retrofit
import retrofit2.adapter.rxjava2.RxJava2CallAdapterFactory
import retrofit2.converter.gson.GsonConverterFactory
import retrofit2.http.*

interface Api {

    companion object {
        fun create(): Api {

            val logging = HttpLoggingInterceptor()
            logging.level = HttpLoggingInterceptor.Level.BODY

            val client = OkHttpClient.Builder()
                .addInterceptor(logging)
                .build()

            val retrofit = Retrofit.Builder()
                .client(client)
                .addCallAdapterFactory(RxJava2CallAdapterFactory.create())
                .addConverterFactory(GsonConverterFactory.create())
                .baseUrl("http://45.128.204.157/")
                .build()

            return retrofit.create(Api::class.java)
        }
    }


    @POST("login")
    fun login(@Body body:LoginRequest): Single<Response<LoginResponse>>

    @POST("/json/obr")
    fun addDocument(@Body body: Document,token:String):Single<Response<CreateResponse>>
}