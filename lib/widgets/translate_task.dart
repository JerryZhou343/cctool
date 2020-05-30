import 'package:cctool/model/translate_task.dart';
import 'package:flutter/material.dart';
import 'package:cctool/model/choice.dart';
import 'package:cctool/common/const.dart';

Widget buildTranslateTask(BuildContext context) {
  var task = TranslateTask();
  task.from = "en";
  task.to = "zh";
  task.merge = false;

  var style = TextStyle(fontSize: 16, decoration: TextDecoration.none);
  return SimpleDialog(
    children: <Widget>[
      Form(child: Column(children: <Widget>[])),
    ],
  );
}
