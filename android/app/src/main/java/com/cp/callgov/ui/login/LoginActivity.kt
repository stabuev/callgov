package com.cp.callgov.ui.login

import android.content.Intent
import androidx.appcompat.app.AppCompatActivity
import android.os.Bundle
import androidx.lifecycle.Observer
import androidx.lifecycle.ViewModelProviders
import com.cp.callgov.DocumentCreationActivity
import com.cp.callgov.R
import com.cp.callgov.ui.main.MainActivity
import com.google.android.material.snackbar.Snackbar
import kotlinx.android.synthetic.main.activity_login.*

class LoginActivity : AppCompatActivity() {

    private lateinit var viewModel:LoginViewModel
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_login)
        viewModel = ViewModelProviders.of(this).get(LoginViewModel::class.java)
        viewModel.getLiveData().observe(this, Observer {
            if(it=="Sucess") startMain()
            else
                showSnackbar(it)
        })
        loginButton.setOnClickListener {
            val login = email.text.toString()
            val password = password.text.toString()
            viewModel.login(login,password)
        }
    }

    private fun showSnackbar(s: String) {
        Snackbar.make(loginContainer,s,Snackbar.LENGTH_SHORT).show()
    }

    private fun startMain() {
        startActivity(Intent(this, MainActivity::class.java))
        finish()
    }

}
