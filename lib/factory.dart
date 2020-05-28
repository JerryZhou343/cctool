
import 'package:flutter/material.dart';
import 'package:cctool/choice.dart';
import 'package:cctool/const.dart';



List<Widget> ConstructWidget(Choice choice, TextStyle textStyle){
  if (choice.title != Setting){
    return <Widget>[
      Icon(choice.icon, size: 128.0, color: textStyle.color),
      Text(choice.title, style: textStyle),
      buildButtonColumn(Icons.call, 'CALL'),
    ];
  }else{
    return <Widget>[
      //Icon(choice.icon, size: 108.0, color: textStyle.color),
      Text(choice.title, style: textStyle),
    ];
  }
}



Column buildButtonColumn(IconData icon, String label) {
  Color color = Color(0xFFFF9000);

  return new Column(
    mainAxisSize: MainAxisSize.min,
    mainAxisAlignment: MainAxisAlignment.center,
    children: [
      new Icon(icon, color: color),
      new Container(
        margin: const EdgeInsets.only(top: 8.0),
        child: new Text(
          label,
          style: new TextStyle(
            fontSize: 12.0,
            fontWeight: FontWeight.w400,
            color: color,
          ),
        ),
      ),
    ],
  );
}
