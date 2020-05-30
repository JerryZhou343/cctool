import 'package:cctool/model/choice.dart';

class Provider {
  static Provider _instance;
  static Choice _choice;

  static Provider getInstance() {
    if (_instance == null) {
      _instance = Provider._internal();
    }
    return _instance;
  }

  Provider._internal();

  SetChoice(Choice choice) {
    _choice = choice;
  }

  Choice GetChoice() {
    return _choice;
  }
}
