package com.cp.callgov

import android.content.Context
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.ImageView
import android.widget.TextView
import androidx.recyclerview.widget.RecyclerView
import com.cp.callgov.model.Document

class ObrsAdapter(private val context: Context, private val obrs: MutableList<Document>) :
    RecyclerView.Adapter<ObrsAdapter.ViewHolder>() {


    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): ViewHolder {
        val view = LayoutInflater.from(parent.context).inflate(R.layout.obr_item, parent, false)
        return ViewHolder(view)
    }

    override fun getItemCount() = obrs.size


    override fun onBindViewHolder(holder: ViewHolder, position: Int) {
        holder.setData(obrs[position],position)
    }

    inner class ViewHolder(view: View) : RecyclerView.ViewHolder(view) {
        private var text: TextView = view.findViewById(R.id.obrTextView)
        private var bottomText: TextView = view.findViewById(R.id.textView)
        private var pos = 0
        var id: String? = null

        init {
            view.setOnClickListener {
                context.startDetail(obrs[pos])
            }
        }

        fun setData(doc: Document,position: Int) {
            this.pos = position
            text.text = doc.title
            bottomText.text = doc.content
        }
    }

}