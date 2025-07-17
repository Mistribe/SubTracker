import 'package:flutter/foundation.dart';
import '../models/label.dart';
import '../repositories/label_repository.dart';

/// Provider for managing labels
///
/// This class provides methods for adding, updating, and retrieving labels.
/// It uses the label repository for data persistence.
class LabelProvider with ChangeNotifier {
  final LabelRepository labelRepository;

  LabelProvider({required this.labelRepository});

  /// Get all labels
  List<Label> get labels => List.unmodifiable(labelRepository.getAll());

  /// Get default labels
  List<Label> get defaultLabels => labelRepository.getDefaultLabels();

  /// Get custom labels
  List<Label> get customLabels => labelRepository.getCustomLabels();

  /// Get a label by ID
  Label? getLabel(String id) => labelRepository.get(id);

  /// Add a new custom label
  Future<Label> addLabel(String name, String color) async {
    final label = await labelRepository.add(name, color);
    notifyListeners();
    return label;
  }

  /// Update an existing label
  Future<void> updateLabel(Label label) async {
    await labelRepository.update(label);
    notifyListeners();
  }

  /// Delete a label
  Future<void> deleteLabel(String id) async {
    await labelRepository.delete(id);
    notifyListeners();
  }
}
