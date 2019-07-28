//
//  LogInViewController.swift
//  CallGov
//
//  Created by Zaur Kasaev on 28/07/2019.
//  Copyright © 2019 callgov. All rights reserved.
//

import UIKit
import Alamofire
import SwiftyJSON

class LogInViewController: UIViewController {

    @IBOutlet weak var login: UITextField!
    @IBOutlet weak var password: UITextField!
    
    var delegate = UIApplication.shared.delegate as! AppDelegate
    
    @IBOutlet weak var progressView: UIActivityIndicatorView!
    override func viewDidLoad() {
        super.viewDidLoad()
        // Do any additional setup after loading the view.
    }

    @IBAction func LoginAction(_ sender: UIButton) {
        if login.text != "" && password.text != ""{
            progressView.startAnimating()
            UIApplication.shared.beginIgnoringInteractionEvents()
            
            let params: [String: Any] = ["login": login.text!,"password": password.text!]
            request("http://45.128.204.157/login", method: .post, parameters: params, encoding: JSONEncoding.default).validate().responseJSON { response in
                
                switch response.result {
                    
                case .success(let value) :
                    
                    let json = JSON(value)
                    
                    self.delegate.token = json["token"].stringValue
                    
                    self.progressView.stopAnimating()
                    UIApplication.shared.endIgnoringInteractionEvents()
                    
                    let storyBoard: UIStoryboard = UIStoryboard(name: "Main", bundle: nil)
                    
                    var ViewController = UIViewController()
                    
                    ViewController = storyBoard.instantiateViewController(withIdentifier: "ViewController") as! ViewController
                   
                    self.present(ViewController, animated: true, completion: nil)
                    
                  
                    
//                    print(json["token"].stringValue)
                    
                case .failure(let error):
                    print(error.localizedDescription)
                }
            }
            
        }
    }
 
    
    

}
