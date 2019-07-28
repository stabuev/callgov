//
//  ViewController.swift
//  CallGov
//
//  Created by Zaur Kasaev on 27/07/2019.
//  Copyright © 2019 callgov. All rights reserved.
//

import UIKit
import Alamofire
import SwiftyJSON

class ViewController: UIViewController {

    @IBOutlet weak var paperCollectionView: UICollectionView!
    
    
    var list : [PaperModel] = []
    
    var delegate = UIApplication.shared.delegate as! AppDelegate
    
    override func viewDidLoad() {
        super.viewDidLoad()
        paperCollectionView.delegate = self
        paperCollectionView.dataSource = self
        print("Token : \(delegate.token)")
        
        renderCollectionView()
    }
    
    func renderCollectionView(){
        
        let height = paperCollectionView.frame.height
        let width = paperCollectionView.frame.width
        
        flowLayout(collectionView: paperCollectionView, size: CGSize(width: width * 0.9, height: height * 0.4), minimumLineSpacing: 4.0, scroll: .vertical)
        
    }
    
    
    private func flowLayout(collectionView : UICollectionView, size : CGSize, minimumLineSpacing : CGFloat, scroll : UICollectionView.ScrollDirection){
        let flowLayout = UICollectionViewFlowLayout()
        flowLayout.itemSize = size
        flowLayout.scrollDirection = scroll
        flowLayout.minimumLineSpacing = minimumLineSpacing
        flowLayout.minimumInteritemSpacing = 0
        flowLayout.headerReferenceSize = CGSize(width: 0, height: 0)
        flowLayout.sectionInset = UIEdgeInsets(top: 0, left: 0, bottom: 0, right: 0)
        collectionView.collectionViewLayout = flowLayout
        collectionView.backgroundColor = UIColor.clear
    }
    
    override func viewWillLayoutSubviews() {
        request("http://45.128.204.157/json/obrlist").validate().responseJSON { (resp) in
            switch resp.result {
            case .success(let value):
                
                let json = JSON(value)
                
                self.list.removeAll()
                for (_,value) in json["obr"].arrayValue.enumerated(){
                    
                    print(value)
                    
                    self.list.append(PaperModel.init(id: value[0].int64Value, title: value[1].stringValue, content: value[2].stringValue, name: value[3].stringValue, public1: value[4].stringValue, state: value[5].stringValue, address: value[6].stringValue, dtreg: value[7].stringValue, dtlast: value[8].stringValue, signatures : value[9].stringValue, userSign : value[10].stringValue))
                }
                self.paperCollectionView.reloadData()
                
            //                print(value)
            case .failure(let error):
                print(error.localizedDescription)
            }
        }
        
        
        
    }

}
extension ViewController : UICollectionViewDelegate, UICollectionViewDataSource{
    
    func collectionView(_ collectionView: UICollectionView, numberOfItemsInSection section: Int) -> Int {
        return list.count
    }
    
    func collectionView(_ collectionView: UICollectionView, cellForItemAt indexPath: IndexPath) -> UICollectionViewCell {
        let cell = collectionView.dequeueReusableCell(withReuseIdentifier: "PaperCell", for: indexPath) as! PaperCollectionViewCell
    
        cell.address.text = "Адресат : \(list[indexPath.item].address)"
        cell.author.text = "Автор :\(list[indexPath.item].name)"
        cell.comments.text = "Комментарии : \(list[indexPath.item].public1)"
        cell.discription.text =  list[indexPath.item].content
        cell.title.text = list[indexPath.item].title
        cell.signatures.text = "Подписей : \(list[indexPath.item].signatures)"
        
        cell.backview.layer.borderColor =  UIColor.black.cgColor
        cell.backview.layer.borderWidth = 0.6
        
        return cell
    }
    
}
