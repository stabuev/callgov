package com.cp.callgov.ui.main

import android.content.Intent
import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import androidx.recyclerview.widget.LinearLayoutManager
import androidx.recyclerview.widget.RecyclerView
import com.cp.callgov.DocumentCreationActivity
import com.cp.callgov.ObrsAdapter
import com.cp.callgov.R
import com.cp.callgov.api.ApiService
import com.cp.callgov.model.Document
import io.reactivex.android.schedulers.AndroidSchedulers
import io.reactivex.rxkotlin.subscribeBy
import io.reactivex.schedulers.Schedulers
import kotlinx.android.synthetic.main.activity_main.*
import kotlinx.android.synthetic.main.content_main.*
import org.json.JSONArray
import org.json.JSONObject

class MainActivity : AppCompatActivity() {

    private val obrs = mutableListOf<Document>()
    private lateinit var adapter:ObrsAdapter
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
        setSupportActionBar(toolbar)

        val d = ApiService.fetchDocuments()
            .subscribeOn(Schedulers.io())
            .observeOn(AndroidSchedulers.mainThread())
            .subscribeBy(
                onSuccess = {
                    val list = JSONObject(it).getJSONArray("obr")
                    for(i in 0 until list.length()){
                        val tmp = list.getJSONArray(i)
                        val document = Document(
                            id = tmp.getInt(0),
                            title = tmp.getString(1),
                            content = tmp.getString(2),
                            state =  tmp.getString(5),
                            address = tmp.getString(6),
                            public = tmp.getInt(4)
                        )
                        obrs.add(document)

                    }
                    val layoutManager = LinearLayoutManager(this, RecyclerView.VERTICAL, false)
                    adapter = ObrsAdapter(this,obrs)
                    recyclerView.adapter = adapter
                    recyclerView.layoutManager = layoutManager
                    adapter.notifyDataSetChanged()
                }
            )
        fab.setOnClickListener {
            startActivity(Intent(this, DocumentCreationActivity::class.java))
        }
    }

}
