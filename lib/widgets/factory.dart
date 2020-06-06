import 'package:flutter/material.dart';
import 'package:cctool/model/choice.dart';
import 'package:cctool/common/const.dart';
import 'package:cctool/widgets/translate_task.dart';

List<Widget> ConstructWidget(Choice choice, TextStyle textStyle) {
  if (choice.title != Setting) {
    return <Widget>[
      Icon(choice.icon, size: 128.0, color: textStyle.color),
      Text(choice.title, style: textStyle),
    ];
  } else {
    return <Widget>[
      Icon(choice.icon, size: 128.0, color: textStyle.color),
      Text(choice.title, style: textStyle),
    ];
  }
}

Widget buildActionButton(BuildContext context, Choice choice) {
  if (choice.title != Setting) {
    return FloatingActionButton(
      onPressed: () {
        // Add your onPressed code here!
        switch (choice.title) {
          case Translate:
            {
              showDialog(
                // 传入 context
                context: context,
                // 构建 Dialog 的视图
                builder: buildTranslateTask,
              );
            }
        }
      },
      child: Icon(Icons.plus_one),
      backgroundColor: Colors.blue,
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
