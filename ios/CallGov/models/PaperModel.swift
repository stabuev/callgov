//
//  PaperModel.swift
//  CallGov
//
//  Created by Zaur Kasaev on 28/07/2019.
//  Copyright © 2019 callgov. All rights reserved.
//

import Foundation

class PaperModel {
    var id : Int64
    var title : String
    var content : String
    var name : String
    var public1 : String
    var state : String
    var address : String
    var dtreg : String
    var dtlast : String
    var signatures : String
    var userSign : String
    
    init(id:Int64,title:String,content:String, name:String, public1:String, state:String, address:String, dtreg:String, dtlast:String, signatures : String, userSign : String) {
        self.id = id
        self.title = title
        self.content = content
        self.name = name
        self.public1 = public1
        self.state = state
        self.address = address
        self.dtreg = dtreg
        self.dtlast = dtlast
        self.signatures = signatures
        self.userSign = userSign
    }
    
}
