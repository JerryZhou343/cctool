import 'package:flutter/material.dart';
import 'file:///E:/GoPathBase/src/github.com/JerryZhou343/cctool/lib/model/choice.dart';
import 'file:///E:/GoPathBase/src/github.com/JerryZhou343/cctool/lib/widgets/factory.dart';

class ChoiceCard extends StatelessWidget {
  const ChoiceCard({Key key, this.choice}) : super(key: key);

  final Choice choice;

  @override
  Widget build(BuildContext context) {
    final TextStyle textStyle = Theme.of(context).textTheme.headline4;
    return Card(
      color: Colors.white,
      child: Center(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.center,
          children: ConstructWidget(this.choice, textStyle),
        ),
      ),
    );
  }
}
