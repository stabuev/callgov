package com.cp.callgov

import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import android.widget.ArrayAdapter
import com.cp.callgov.api.ApiService
import com.cp.callgov.model.Document
import io.reactivex.android.schedulers.AndroidSchedulers
import io.reactivex.rxkotlin.subscribeBy
import io.reactivex.schedulers.Schedulers
import kotlinx.android.synthetic.main.activity_document_creation.*

class DocumentCreationActivity : AppCompatActivity() {

    lateinit var d:Document
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_document_creation)
        d = intent.getSerializableExtra("obr") as Document

        val typeAdapter =  ArrayAdapter.createFromResource(this, R.array.pp, android.R.layout.simple_spinner_item);
        typeAdapter.setDropDownViewResource(android.R.layout.simple_spinner_dropdown_item);

        val readyAdapter =  ArrayAdapter.createFromResource(this, R.array.work_status, android.R.layout.simple_spinner_item);
        readyAdapter.setDropDownViewResource(android.R.layout.simple_spinner_dropdown_item);
        typeSpinner.adapter = typeAdapter
        sostSpinner.adapter = readyAdapter
        button.setOnClickListener {
            val d = Document(
                title = titleInput.text.toString(),
                address = address.text.toString(),
                content = content.text.toString().replace("\n"," "),
                public = if(typeSpinner.selectedItem.toString()=="Публичное") 1; else 0,
                state = "draft"
            )
            ApiService.createDocument(d)
                .subscribeOn(Schedulers.io())
                .observeOn(AndroidSchedulers.mainThread())
                .subscribeBy {
                    finish()
                }
        }
    }
}
