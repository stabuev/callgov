//
//  PaperCollectionViewCell.swift
//  CallGov
//
//  Created by Zaur Kasaev on 28/07/2019.
//  Copyright © 2019 callgov. All rights reserved.
//

import UIKit

class PaperCollectionViewCell: UICollectionViewCell {
    
    @IBOutlet weak var title: UILabel!
    @IBOutlet weak var discription: UITextView!
    @IBOutlet weak var author: UILabel!
    @IBOutlet weak var backview: GrView!
    @IBOutlet weak var address: UILabel!
    @IBOutlet weak var signatures: UILabel!
    @IBOutlet weak var comments: UILabel!
    
}
