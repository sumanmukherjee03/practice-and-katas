# @param {Integer[][]} points
# @param {Integer} k
# @return {Integer[][]}
def k_closest(points, k)
    dists = []
    h = {}
    points.each do |p|
        d = Math.sqrt(p[0]**2 + p[1]**2).round(2)
        dists << d
        if h.has_key?(d)
            h[d] << p
        else
            h[d] = [p]
        end
    end
    res = dists.uniq.sort[0...k].inject([]) {|a, d| a.concat(h[d])}[0...k]
    return res
end
