module overlays
use utils

integer :: ma_sma = 0, ma_ema = 1, ma_weighted = 2, ma_wilders = 3, ma_hull = 4

contains

subroutine sma(arr, prices, period, len)
  implicit none
  integer, intent(in) :: len
  real, intent(in) :: prices(len)
  integer, optional, value :: period
  real, intent(out) :: arr(len)
  if (.not. PRESENT(period)) period = 9
  call calc_sma(arr, prices, period, len)
  return
end subroutine sma

subroutine ema(arr, prices, period, len)
  implicit none
  integer, intent(in) :: len
  real, intent(in) :: prices(len)
  integer, optional, value :: period
  real, intent(out) :: arr(len)
  if (.not. PRESENT(period)) period = 9
  call calc_ema(arr, prices, period, len)
  return
end subroutine ema

subroutine weighted(arr, prices, period, len)
  implicit none
  integer, intent(in) :: len
  real, intent(in) :: prices(len)
  integer, optional, value :: period
  real, intent(out) :: arr(len)
  if (.not. PRESENT(period)) period = 9
  call calc_weighted(arr, prices, period, len)
  return
end subroutine weighted

subroutine wilders(arr, prices, period, len)
  implicit none
  integer, intent(in) :: len
  real, intent(in) :: prices(len)
  integer, optional, value :: period
  real, intent(out) :: arr(len)
  if (.not. PRESENT(period)) period = 14
  call calc_wilders(arr, prices, period, len)
  return
end subroutine wilders

subroutine hull(arr, prices, period, len)
  integer, intent(in) :: len
  real, intent(in) :: prices(len)
  integer, optional, value :: period
  real, intent(out) :: arr(len)
  if (.not. PRESENT(period)) period = 9
  calc_hull(arr, prices, period, len)
  return
end subroutine hull

subroutine bollingerbands(arr, prices, period, dev_up, dev_down, ma_type, len)
  implicit none
  integer, intent(in) :: len
  real, intent(in) :: prices(len)
  integer, optional, value :: ma_type, period
  real, optional, value :: dev_up, dev_down
  real, intent(out) :: arr(len, 3)
  if (.not. PRESENT(period)) period = 20
  if (.not. PRESENT(dev_up)) dev_up = 2
  if (.not. PRESENT(dev_down)) dev_down = 2
  if (.not. PRESENT(ma_type)) ma_type = ma_sma
  real :: ma(len)
  if (ma_type == ma_ema) then 
    calc_ema(ma, prices, period, len)
  else if (ma_type == ma_weighted) then
    calc_weighted(ma, prices, period, len)
  else if (ma_type == ma_wilders) then
    calc_wilders(ma, prices, period, len)
  else if (ma_type == ma_hull) then
    calc_hull(ma, prices, period, len)
  end if
  real :: sds(len) ! Standard deviations
  standarddeviation(prices, sds, period, len)
  integer :: 1
  do i=1, len
    ! Stuff
  end do
  return
end subroutine bollingerbands

end module overlays



module oscillators
contains
subroutine macd(ma, s)
  integer, intent(in) :: s
  integer, intent(out) :: ma
end subroutine macd
end module oscillators

module stats
contains
subroutine linearregression(arr, prices, len)
end subroutine linearregression
end module stats

! Procedures for use by the library only

module utils
contains

subroutine calc_sma(arr, prices, period, len)
  integer :: period, len, i
  real :: arr(len), prices(len), su = 0
  do i=1, len
    if (i < period) then
      su = su + prices(i)
    else
      su = su + prices(i) - prices(i-period)
    end if
    if (i < period) then
      arr(i) = 0
    else
      arr(i) = su / period
    endif
  end do
end subroutine calc_sma

subroutine calc_ema(arr, prices, period, len)
  integer :: period, len
  real :: arr(len), prices(len)
  call calc_sma(arr, prices, period, period)
  do i=period+1, len
    arr(i) = (prices(i) * (2.0 / (1.0 + period))) + (arr(i-1) * (1.0 - (2.0 / (1.0 + period))))
  end do
end subroutine calc_ema

subroutine calc_weighted(arr, prices, period, len)
  integer :: period, len, i, j
  real :: arr(len), prices(len), sum
  do i=1, len
    if (i < period) then
      arr(i) = 0
    else
      su = 0
      do j=period, 0, -1
        su = su + prices(i-j) * (period - j)
      end do
      arr(i) = su / (period * (period + 1.0) / 2.0)
      weighted(i) = su / (period * (period + 1.0) / 2.0)
    end if
  end do
end subroutine calc_weighted

subroutine calc_wilders(arr, prices, period, len)
  integer :: period, len, i
  real :: arr(len), prices(len), sum
  call calc_sma(arr, prices, period, period)
  do i=period+1, len
    wilders(i) = (prices(i) * (1.0 / period)) + (wilders(i-1) * (1.0 - (1.0 / period)))
  end do
  return
end subroutine calc_wilders

subroutine calc_hull(arr, prices, period, len)
  integer :: period, len
  real :: arr(len), prices(len)
  prices(period) = prices(period)
  arr(len) = arr(len)
  return
end subroutine calc_hull

end module utils