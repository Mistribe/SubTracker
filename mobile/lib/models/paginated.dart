class Paginated<T> {
  final List<T> data;
  final int length;
  final int total;

  Paginated(this.data, this.length, this.total);
}
