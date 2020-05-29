import 'package:flutter/material.dart';
import 'file:///E:/GoPathBase/src/github.com/JerryZhou343/cctool/lib/model/choice.dart';
import 'file:///E:/GoPathBase/src/github.com/JerryZhou343/cctool/lib/common/const.dart';

List<Widget> ConstructWidget(Choice choice, TextStyle textStyle) {
  if (choice.title != Setting) {
    return <Widget>[
      Icon(choice.icon, size: 128.0, color: textStyle.color),
      Text(choice.title, style: textStyle),
      Positioned(top:1000,left: 1000,child:buildActionButton(choice))
    ];
  } else {
    return <Widget>[
      Icon(choice.icon, size: 128.0, color: textStyle.color),
      Text(choice.title, style: textStyle),
      buildActionButton(choice)
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

Widget buildActionButton(Choice choice) {
  if (choice.title != Setting) {
    return FloatingActionButton(
      onPressed: () {
        // Add your onPressed code here!
      },
      child: Icon(Icons.navigation),
      backgroundColor: Colors.green,
    );
  } else {
    return FloatingActionButton.extended(
      onPressed: () {
        // Add your onPressed code here!
      },
      label: Text('Save'),
      icon: Icon(Icons.thumb_up),
      backgroundColor: Colors.pink,
    );
  }
}
