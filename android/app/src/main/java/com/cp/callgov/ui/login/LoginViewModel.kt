package com.cp.callgov.ui.login

import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import com.cp.callgov.api.Api
import com.cp.callgov.api.ApiService
import io.reactivex.android.schedulers.AndroidSchedulers
import io.reactivex.disposables.CompositeDisposable
import io.reactivex.rxkotlin.subscribeBy
import io.reactivex.schedulers.Schedulers

class LoginViewModel : ViewModel() {
    private val disposables = CompositeDisposable()
    private val loginLiveData = MutableLiveData<String>()
    fun login(login: String, password: String) {
        if ((login.isEmpty() || password.isEmpty())) loginLiveData.postValue("Неправильный логин или пароль")
        else {
            disposables.add(ApiService.login(login, password)
                .subscribeOn(Schedulers.io())
                .observeOn(AndroidSchedulers.mainThread())
                .subscribeBy(
                    onSuccess = {
                        val body = it.body()!!.token
                        if (body.isNullOrEmpty())
                            loginLiveData.postValue("Неправильный логин или пароль")
                        else loginLiveData.postValue("Sucess")
                        ApiService.token = body
                    },
                    onError = {
                        loginLiveData.postValue("Ошибка")
                    }
                ))
        }
    }

    fun getLiveData() = loginLiveData

    override fun onCleared() {
        disposables.dispose()
    }
}